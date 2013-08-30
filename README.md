Goson
=====
Goson is a small simple DSL for generating json from your Go datatypes. Supports both structs and maps as well as function/method values.

Installation
------------
Make sure to first install Go and setup all the necessary environment variables. Visit http://golang.org/doc/install for more info.
After that simply `go get github.com/emilsjolander/goson`.

Getting started
---------------
This is a short demonstration of the library in use. The next sections will go into detail about the API and the templating syntax.

All the code below can be found in the sample folder of the project.

Start with setting up a small Go program, something like this:
```go
package main

import (
	"fmt"
	"github.com/emilsjolander/goson"
)

type Repo struct {
	Name  string
	Url   string
	Stars int
	Forks int
}

type User struct {
	Name  string
	Repos []Repo
}

func main() {

	user := &User{
		Name: "Emil Sjölander",
		Repos: []Repo{
			Repo{
				Name:  "goson",
				Url:   "https://github.com/emilsjolander/goson",
				Stars: 0,
				Forks: 0,
			},
			Repo{
				Name:  "StickyListHeaders",
				Url:   "https://github.com/emilsjolander/StickyListHeaders",
				Stars: 722,
				Forks: 197,
			},
			Repo{
				Name:  "android-FlipView",
				Url:   "https://github.com/emilsjolander/android-FlipView",
				Stars: 157,
				Forks: 47,
			},
		},
	}

	result, err := goson.Render("user", goson.Args{"User": user})

	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
}
```
The first thing we do in the above code is to import goson as well as the fmt packages. We also define 2 data types, User and Repo. These are the data types that we want to format into json. In the main function of our small sample we create an instance of User containing 3 instances of Repo. After this come the only public function of this library, the `Render()` function. `Render` takes two arguments, first the name of the template to render excluding the templates file type which should be `.goson`. The second argument to render is a map of argument that the tempate can make use of.

Let's take a look at user.goson to see how we define our json structure.
```text
user: {
	name: User.Name
	repos: Repo in User.Repos {
		name: Repo.Name
		url: Repo.Url
		stars: Repo.Stars
		forks: Repo.Forks
	}
}
```
The Above template starts by wrapping the fields within a "user" json object, next it writes the name of the user and than itterates through the repos printing each repos name, url, stars and forks. 
The resulting json is the following:
```json
{
    "user": {
        "name": "Emil Sjölander",
        "repos": [
            {
                "name": "goson",
                "url": "https://github.com/emilsjolander/goson",
                "stars": 0,
                "forks": 0
            },
            {
                "name": "StickyListHeaders",
                "url": "https://github.com/emilsjolander/StickyListHeaders",
                "stars": 722,
                "forks": 197
            },
            {
                "name": "android-FlipView",
                "url": "https://github.com/emilsjolander/android-FlipView",
                "stars": 157,
                "forks": 47
            }
        ]
    }
}
```
As you can see, the result is automatically wrapped inside a json object. This is to follow standard restful response formats.

Why?
----
You might ask why should I use this over Go's built in encoding/json package? That's a fair question, you might not have any need for goson. If you are building an API server you most likely have use for goson though. The json marshaler in encoding/json in both quick and fairly easy to use but it is not flexible or secure. By not being secure I mean that it is easy to leak private field when `encoding/json` uses a opt-out strategy for json fields. This is where goson comes into play!

Goson lets you render the same data type into different json output depending on the situation. You might have both a public and private API, in this case you could have a `templates/private/user.goson` and a `templates/public/user.goson` template, the public template might skip some internal fields as an auth token or perhaps the id of the user. One other time where goson is very useful is in the above sample, to save space I might just want to render the url of a repo when the user of the API GETs /user/1 but when they GET /user/1/repo/1 I will render all the info attached to the repo.

Another reason to use goson is that it separates the view layer(json in this case) from the model layer. Defining the json keys within the model is against any good MVC design and should be avoided when possible.

API
---
As I hinted at during the getting started part of this readme, the API is very small. It consists of only one function and that is
```go
goson.Render(template string, args Args)
```
The template parameter should be the relative filepath to the template. So if you are executing main.go and your template is inside the templates folder you will want to pass `"templates/my_template"` to `Render()`. This will render the your data with the my_template.goson template which is located inside the template directory.

Args is just an alias for map[string]interface{} and accepts almost anything as an argument. Complex numbers and channels are the two common data types not currently supported.

Syntax
------
Goson is a fairly powerful templating language with support for everything you could want (Open a pull request if I've missed anything).

define a json key the following way.
```text
key_name:
```

After a json key definition there are multiple options that can follow.
A `string`, `int`, `float` or `bool` literal for example.
```text
my_int: 5
my_float: 4.3
my_bool: true
my_string: "Hello world"
```

A json key can also be followed by a json object literal as below.
```text
my_object: {
	nested_string: "Hi there!"
}
```

When defining a object literal you have the possibility to add a alias to a variable for the scope of the object.
```text
my_object: Object.NestedObject.NestedObject as o {
	key: o.value
}
```
`Object.NestedObject.NestedObject` can be either a `struct`, a `*struct` or a `map[string]`.

Looping over a collection works much in the same way.
```text
my_array: o in Object.MyCollection {
	key: o.value
}
```
`Object.MyCollection` can be an instance of either an `Array`, a `Slice` or a `goson.Collection`.

A very important feature for larger applications is to write modular code. The include statement makes this possible.
```text
include(template_name, MyObject.NestedObject)
```
The above code will look for a template named "template_name.goson" in the same folder as the template that included it. `MyObject.NestedObject` is sent as a parameter to partial template. If `MyObject.NestedObject` has a field called `MyField` the partial template will refer to it via `MyField` and not `MyObject.NestedObject.MyField`. The argument (in this case `MyObject.NestedObject`) that is sent to the partial can be either a `struct`, a `*struct` or a `map[string]`.

Comments are also supported, both single and multi line comments. They follow the standard go syntax.
```text
my_object: {
	//TODO add some property
	my_key: MyObject.Key
}
```

Contributing
------------

Pull requests and issues are very welcome!

If you want a fix to happen sooner than later, I suggest that you make a pull request.

Preferably send pull requests early, even before you are done with the feature/fix/enhancement. This way we can discuss and help each other out :)


License
-------

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
