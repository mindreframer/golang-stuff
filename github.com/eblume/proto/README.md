proto
=====

`proto` gives Go operations like `Map`, `Reduce`, `Filter`, `De/Multiplex`, etc.
without sacrificing idiomatic harmony or speed. It also introduces a convenience
type for these functions, `Proto`, which is a stand-in for the empty interface
(interface{}), which is used to box values being sent to these operations.

Documentation
-------------

Please see documentation.{txt,html} for the automatically generated
documentation - or better yet, just run:

    godoc github.com/eblume/proto | less

That's probably a better idea since there's a decent chance the documentation
might be lagging behind the current code base, since it has to be run manually
(at this moment).

You can also take a look at the *_test.go files for an even better look in to
how to use Proto. I will make one disclaimer, which is that code written with
Proto has some unavoidable boilerplate in the form of casting to/from the Proto
type - this boilerplate is annoying but is much less obvious and significant
with larger code bases that use Proto-style channels in chains.

Examples
--------

Double every integer in a slice:

    inputs := []Proto{0, 1, 2, 3, 4, 5, 6}
    sent := Send(inputs)
    doubler := func(a Proto) Proto {
        return a.(int) * 2
    }
    mapped := Map(doubler, sent)
    doubled := Gather(mapped)

Double every integer, chained:

    doubled := Gather(Map(func(a Proto) Proto {
        return a.(int) * 2
    }, Send([]Proto{0, 1, 2, 3, 4, 5, 6})))

License
-------

Please see COPYING for more details on the licensing of this software.
