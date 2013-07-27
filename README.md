This is a tutorial on Go's database/sql package (http://golang.org/pkg/database/sql/).
The package's documentation tells you what everything does, but it
doesn't tell you how to use the package. We find ourselves wishing for a quick-reference and a "getting started"
orientation. This repo is an attempt to provide that. Contributions are welcome.

Go's database/sql Package
=========================

The idiomatic way to use a SQL, or SQL-like, database in Go is through the `database/sql`
package. It provides a lightweight interface to a row-oriented database. This documentation
is a reference for the most common aspects of how to use it.

The first thing to do is import the `database/sql` package, and a driver package. You generally shouldn't use the driver package directly, although some drivers encourage you to do so. (In our opinion, it's usually a bad idea.) Instead, your code should only refer to `database/sql`. This helps avoid making your code dependent on the driver, so that you can change the underlying driver (and thus the database you're accessing) without changing your code. It also forces you to use the Go idioms instead of ad-hoc idioms that a particular driver author may have provided.

In this documentation, we'll use the excellent MySQL drivers at https://github.com/go-sql-driver/mysql for examples.

Add the following to the top of your Go source file:

```go
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
```

Notice that we're loading the driver anonymously, aliasing its package qualifier to `_` so none of its exported names are visible to our code. Under the hood, the driver registers itself as being available to the `database/sql` package, but in general nothing else happens.

Now you're ready to access a database.

Accessing the Database
======================

Now that you've loaded the driver package, you're ready to create a database object, a `sql.DB`. The first thing you should know is that **a `sql.DB` isn't a database connection**. It also doesn't map to any particular database software's notion of a "database" or "schema." It's an abstraction of the interface and existence of a database, which might be a local file, accessed through a network connection, in-memory and in-process, or what have you.

The `sql.DB` performs some important tasks for you behind the scenes:

* It opens and closes connections to the actual underlying database, via the driver.
* It manages a pool of connections as needed.

These "connections" may be file handles, sockets, network connections, or other ways to access the database. The `sql.DB` abstraction is designed to keep you from worrying about how to manage concurrent access to the underlying datastore. A connection is marked in-use when you use it to perform a task, and then returned to the available pool when it's not in use anymore. One consequence of this is that **if you fail to release connections back to the pool, you can cause `db.SQL` to open a lot of connections**, potentially running out of resources (too many connections, too many open file handles, lack of available network ports, etc). We'll discuss more about this later.

To create a `sql.DB`, you use `sql.Open()`. This returns a `*sql.DB`:

```go
func main() {
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
```

In the example shown, we're illustrating several things:

1. The first argument to `sql.Open` is the driver name. This is the string that the driver used to register itself with `database/sql`, and is conventionally the same as the package name to avoid confusion.
2. The second argument is a driver-specific syntax that tells the driver how to access the underlying datastore. In this example, we're connecting to the "hello" database inside a local MySQL server instance.
3. You should (almost) always check and handle errors returned from all `database/sql` operations.
4. It is idiomatic to `defer db.Close()` if the `sql.DB` should not have a lifetime beyond the scope of the function.

Perhaps counter-intuitively, `sql.Open()` **does not establish any connections to the database**, nor does it validate driver connection parameters. Instead, it simply prepares the database abstraction for later use. The first actual connection to the underlying datastore will be established lazily, when it's needed for the first time. If you want to check right away that the database is available and accessible (for example, check that you can establish a network connection and log in), use `db.Ping()` to do that, and remember to check for errors:

```go
	err = db.Ping()
	if err != nil {
		// do something here
	}
```

Although it's idiomatic to `Close()` the database when you're finished with it, **the `sql.DB` object is designed to be long-lived.** Don't `Open()` and `Close()` databases frequently. Instead, create **one** `sql.DB` object for each distinct datastore you need to access, and keep it until the program is done accessing that datastore. Pass it around as needed, or make it available somehow globally, but keep it open. And don't `Open()` and `Close()` from a short-lived function. Instead, pass the `sql.DB` into that short-lived function as an argument.

If you don't treat the `sql.DB` as a long-lived object, you could experience problems such as poor reuse and sharing of connections, running out of available network resources, sporadic failures due to a lot of TCP connections remaining in TIME_WAIT status, and so on.

Common Database Operations
==========================

Now that you've loaded the driver and opened the `sql.DB`, it's time to use it. There are several idiomatic operations against the datastore:

1. Execute a query that returns rows.
1. Execute a query that returns a single row. There is a shortcut for this special case.
1. Prepare a statement for repeated use, execute it multiple times, and destroy it.
1. Execute a statement in a once-off fashion, without preparing it for repeated use.
1. Modify data and check for the results.
1. Perform transaction-related operations; not discussed at this time.

You should almost always capture and check errors from all functions that return them. There are a few special cases that we'll discuss later where it doesn't make sense to do this.

Go's `database/sql` function names are significant. **If a function name includes `Query`, it is designed to ask a question of the database, and should return a set of rows**, even if it's empty. Statements that don't return rows should not use `Query` functions, for reasons we'll also discuss later.

Fetching Data from the Database
===============================

Let's take a look at an example of how to query the database, working with results. We'll query the `users` table for a user whose `id` is 1, and print out the user's `id` and `name`:

```go
	var (
		id int
		name string
	)
	rows, err := db.Query("select id, name from users where id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
```

Here's what's happening in the above code:

1. We're using `db.Query()` to send the query to the database. We check the error, as usual.
2. We defer `rows.Close()`. This is very important; more on that in a moment.
3. We iterate over the rows with `rows.Next()`.
4. We read the columns in each row into variables with `rows.Scan()`.
5. We check for errors after we're done iterating over the rows.

A couple parts of this are easy to get wrong, and can have bad consequences.

First, as long as there's an open result set (represented by `rows`), the underlying connection is busy and can't be used for any other query. That means it's not available in the connection pool. If you iterate over all of the rows with `rows.Next()`, eventually you'll read the last row, and `rows.Next()` will encounter an internal EOF error and call `rows.Close()` for you. But if for any reason you exit that loop -- an error, an early return, or so on -- then the `rows` doesn't get closed, and the connection remains open. This is an easy way to run out of resources. This is why **you should always `defer rows.Close()`**, even if you also call it explicitly at the end of the loop, which isn't a bad idea. `rows.Close()` is a harmless no-op if it's already closed, so you can call it multiple times. Notice, however, that we check the error first, and only do `rows.Close()` if there isn't an error, in order to avoid a runtime panic.

Second, you should always check for an error at the end of the `for rows.Next()` loop. If there's an error during the loop, you need to know about it. Don't just assume that the loop iterates until you've processed all the rows.

The error returned by `rows.Close()` is the only exception to the general rule that it's best to capture and check for errors in all database operations. If `rows.Close()` throws an error, it's unclear what is the right thing to do. Logging the error message or panicing might be the only sensible thing to do, and if that's not sensible, then perhaps you should just ignore the error.

Assigning Results to Variables
==============================

In the previous section you already saw the idiom for assigning results to variables, a row at a time, with `rows.Scan()`. This is pretty much the only way to do it in Go. You can't get a row as a map, for example. That's because everything is strongly typed. You need to create variables of the correct type and pass pointers to them, as shown.

There are two special cases: nullable columns, and variable numbers of columns, that are a little harder to handle.

Nullable columns are annoying and lead to a lot of ugly code. If you can, avoid them. If not, then you'll need to use special types from the `database/sql` package to handle them. There are types for nullable booleans, strings, integers, and floats. Here's how you use them:

```go
for rows.Next() {
	var s sql.NullString
	err := rows.Scan(&s)
	// check err
	if s.Valid {
	   // use s.String
	} else {
	   // NULL value
	}
}
```

Limitations of the nullable types, and reasons to avoid nullable columns in case you need more convincing:

1. There's no `sql.NullUint64` or `sql.NullYourFavoriteType`.
1. Nullability can be tricky, and not future-proof. If you think something won't be null, but you're wrong, your program will crash, perhaps rarely enough that you won't catch errors before you ship them.
1. One of the nice things about Go is having a useful default zero-value for every variable. This isn't the way nullable things work.

The other special case is assigning a variable number of columns into variables. The `rows.Scan()` function accepts a variable number of `interface{}`, and you have to pass the correct number of arguments. If you don't know the columns or their types, you should use `sql.RawBytes`:

```go
cols, err := rows.Columns()				// Get the column names; remember to check err
vals := make([]sql.RawBytes, len(cols)) // Allocate enough values
ints := make([]interface{}, len(cols)) 	// Make a slice of []interface{}
for i := range ints {
	vals[i] = &ints[i] // Copy references into the slice
}
for rows.Next() {
	err := rows.Scan(vals...)
	// Now you can check each element of vals for nil-ness,
	// and you can use type introspection and type assertions
	// to fetch the column into a typed variable.
}
```

If you know the possible sets of columns and their types, it can be a little less annoying, though still not great. In that case, you simply need to examine `rows.Columns()`, which returns an array of column names.

Preparing Queries
=================

You should, in general, always prepare queries to be used multiple times. The result of preparing the query is a prepared statement, which can have `?` placeholders for parameters that you'll provide when you execute the statement. This is much better than concatenating strings, for all the usual reasons (avoiding SQL injection attacks, for example).

```go
stmt, err := db.Prepare("select id, name from users where id = ?")
if err != nil {
	log.Fatal(err)
}
rows, err := stmt.Query(1)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()
for rows.Next() {
	// ...
}
```

Under the hood, `db.Query()` actually prepares, executes, and closes a prepared statement. That's three round-trips to the database. If you're not careful, you can triple the number of database interactions your application makes! Some drivers can avoid this in specific cases with an addition to `database/sql` in Go 1.1, but not all drivers are smart enough to do that. Caveat Emptor.

Single-Row Queries
==================

If a query returns at most one row, you can use a shortcut around some of the lengthy boilerplate code:

```go
var name string
err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)
if err != nil {
	log.Fatal(err)
}
fmt.Println(name)
```

Errors from the query are deferred until `Scan()` is called, and then are returned from that. You can also call `QueryRow()` on a prepared statement:

```go
stmt, err := db.Prepare("select id, name from users where id = ?")
if err != nil {
	log.Fatal(err)
}
var name string
err = stmt.QueryRow(1).Scan(&name)
if err != nil {
	log.Fatal(err)
}
fmt.Println(name)
```

Statements that Modify Data
===========================

As mentioned previously, you should only use `Query` functions to execute queries -- statements that return rows. Use `Exec`, preferably with a prepared statement, to accomplish an INSERT, UPDATE, DELETE, or other statement that doesn't return rows. The following example shows how to insert a row:

```go
stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
if err != nil {
	log.Fatal(err)
}
res, err := stmt.Exec("Dolly")
if err != nil {
	log.Fatal(err)
}
lastId, err := res.LastInsertId()
if err != nil {
	log.Fatal(err)
}
rowCnt, err := res.RowsAffected()
if err != nil {
	log.Fatal(err)
}
log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
```

Executing the statement produces a `sql.Result` that gives access to statement metadata: the last inserted ID and the number of rows affected.

What if you don't care about the result? What if you just want to execute a statement and check if there were any errors, but ignore the result? Wouldn't the following two statements do the same thing?

```go
_, err := db.Exec("DELETE FROM users")  // OK
_, err := db.Query("DELETE FROM users") // BAD
```

The answer is no. They do **not** do the same thing, and **you should never use `Query()` like this.** The `Query()` will return a `sql.Rows`, which will not be released until it's garbage collected, which can be a long time. During that time, it will continue to hold open the underlying connection, and this anti-pattern is therefore a good way to run out of resources (too many connections, for example).

Surprises, Antipatterns and Limitations
=======================================

We've documented several surprises and antipatterns throughout this tutorial, so please refer back to them if you didn't read them already:

* Opening and closing databases can cause exhaustion of resources.
* Failing to use `rows.Close()` can cause exhaustion of resources.
* Using `Query()` for a statement that doesn't return rows is a bad idea.
* Failing to use prepared statements can lead to a lot of extra database activity.
* Nulls cause annoying problems, which may show up only in production.

There are also a couple of limitations in the `database/sql` package. The interface doesn't give you all-encompassing access to what's happening under the hood. For example, you don't have much control over the pool of connections.

Another limitation, which can be a surprise, is that you can't pass big unsigned integers as parameters to statements if their high bit is set:

```go
_, err := db.Exec("INSERT INTO users(id) VALUES", math.MaxUint64)
```

This will throw an error. Be careful if you use `uint64` values, as they may start out small and work without error, but increment over time and start throwing errors.

Conclusion
==========

We hope you've found this tutorial helpful. Please send pull requests with any contributions! We'd especially appreciate help with missing material, such as how to work with transactions.

This work is licensed under a Creative Commons Attribution-ShareAlike 3.0 Unported License.
