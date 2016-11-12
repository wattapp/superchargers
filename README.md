# Superchargers.io

![](https://travis-ci.org/wattapp/superchargers.svg?branch=master)

The GraphQL API for finding Tesla Superchargers, destination chargers, stores, and service centers.

## Setup

To run Superchargers.io locally just pull down the project:

```sh
go get -u github.com/wattapp/superchargers
```

Run `script/setup`:

```sh
$ script/setup
Checking for dependencies
Fetching required packages
Superchargers.io is ready to go, run:
script/server
```

## What is Superchargers.io?

Superchargers.io is a GraphQL API to progmatically find Tesla Stores, Superchargers, Destination Chargers, and Service centers.

You can filter locations by specific countries, regions, [bounding box](http://wiki.openstreetmap.org/wiki/Bounding_Box), or whether or not the location is open for use.

## What is GraphQL?

Originally developed by Facebook, [GraphQL](http://graphql.org/) is described as:

> GraphQL is a query language for APIs and a runtime for fulfilling those queries with your existing data. GraphQL provides a complete and understandable description of the data in your API, gives clients the power to ask for exactly what they need and nothing more, makes it easier to evolve APIs over time, and enables powerful developer tools.

## I've never written a GraphQL query before, do you have any examples?

Why of course!

- [Superchargers in the Europe region](https://www.superchargers.io/?query=%7B%0A%20%20locations(region%3A%20EUROPE%2C%20type%3A%20SUPERCHARGER)%20%7B%0A%20%20%20%20edges%20%7B%0A%20%20%20%20%20%20node%20%7B%0A%20%20%20%20%20%20%20%20id%0A%20%20%20%20%20%20%20%20address%0A%20%20%20%20%20%20%20%20latitude%0A%20%20%20%20%20%20%20%20longitude%0A%20%20%20%20%20%20%20%20title%0A%20%20%20%20%20%20%20%20locationType%0A%20%20%20%20%20%20%20%20openSoon%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D)
- [Locations in a given bounding box](https://www.superchargers.io/?query=%7B%0A%20%20locations(boundingBox%3A%20%5B42.02238033615207%2C%20-76.4456118%2C%2038.88848331958911%2C%20-83.4768618%5D)%20%7B%0A%20%20%20%20edges%20%7B%0A%20%20%20%20%20%20node%20%7B%0A%20%20%20%20%20%20%20%20id%0A%20%20%20%20%20%20%20%20address%0A%20%20%20%20%20%20%20%20latitude%0A%20%20%20%20%20%20%20%20longitude%0A%20%20%20%20%20%20%20%20title%0A%20%20%20%20%20%20%20%20locationType%0A%20%20%20%20%20%20%20%20openSoon%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D)
- [Locations in the Asia Pacific region](https://www.superchargers.io/?query=%7B%0A%20%20locations(region%3A%20ASIA_PACIFIC)%20%7B%0A%20%20%20%20edges%20%7B%0A%20%20%20%20%20%20node%20%7B%0A%20%20%20%20%20%20%20%20id%0A%20%20%20%20%20%20%20%20address%0A%20%20%20%20%20%20%20%20latitude%0A%20%20%20%20%20%20%20%20longitude%0A%20%20%20%20%20%20%20%20title%0A%20%20%20%20%20%20%20%20locationType%0A%20%20%20%20%20%20%20%20openSoon%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D)
- [Locations that are opening soon](https://www.superchargers.io/?query=%7B%0A%20%20locations(openSoon%3A%20true)%20%7B%0A%20%20%20%20edges%20%7B%0A%20%20%20%20%20%20node%20%7B%0A%20%20%20%20%20%20%20%20id%0A%20%20%20%20%20%20%20%20address%0A%20%20%20%20%20%20%20%20latitude%0A%20%20%20%20%20%20%20%20longitude%0A%20%20%20%20%20%20%20%20title%0A%20%20%20%20%20%20%20%20locationType%0A%20%20%20%20%20%20%20%20openSoon%0A%20%20%20%20%20%20%7D%0A%20%20%20%20%7D%0A%20%20%7D%0A%7D)

In the [GraphiQL interactive query editor](https://www.superchargers.io/graphiql) provided by Superchargers.io, you can use autocomplete and the provided documentation to help craft your queries. If you need any help please open an issue on the [GitHub issue tracker](/wattapp/superchargers/issues) or tweet at [@garrettb](https://twitter.com/garrettb) on Twitter.

## How often are the locations updated?

The locations are updated daily at `00:00 UTC`

## Are you affiliated with Tesla?

Superchargers.io is in no way affiliated with Tesla, Inc. Site names, logos, and images are copyright of Tesla, Inc.

## Something is broken, where can I report it?

You can report issues to the [GitHub issue tracker](/wattapp/superchargers/issues) or directly to [@garrettb](https://twitter.com/garrettb) on Twitter.

## How can I support Superchargers.io?

- Develop a client that consumes Superchargers.io in your favorite language and share it with the repository
- Report any issues to the [GitHub issue tracker](/wattapp/superchargers/issues)
- Tell your developer friends about the service
- Even if you just found it useful or informative, tweet at [@garrettb](https://twitter.com/garrettb) to let me know

## Copyright

Copyright © 2016 Garrett Bjerkhoel. See [MIT-LICENSE](/wattapp/superchargers/blob/master/MIT-LICENSE) for details.
