# CheckMyResult

    > A helper program that scraps and emails results from Anna University servers during high traffic concurrently.

   #### older version of this project: [https://github.com/AravindVasudev/CheckMyResult_OLD]()

## What does this do?

This program takes a JSON array of student's register number with their email ID, fetches their result and emails them. This is useful since the servers are usually overloaded during when the results come out and this program uses a retry function when the request results in an error. This works concurrently and hence is efficient when fetching for multiple students.

## Dependencies

   * [github.com/PuerkitoBio/goquery]()
   * [github.com/Rican7/retry]()

## Installation

1. Install Golang

2. Clone this repository

```
    $ git clone https://github.com/AravindVasudev/CheckMyResult.git
    $ cd CheckMyResult
```

3. Install all dependencies

```
    $ go get ./...
```

4. Build the project

```
    $ go build .
```

5. Create `email_smtp.json` with your stmp server details

```
    $ echo "{\"emailID\": \"email@example.com\",\"password\": \"password\",\"server\": \"smtp.gmail.com\"}" > email_smtp.json
```

6. Create `students.json` with all the student details

```
    $ echo "[{\"registerNumber\": \"123456789\", \"emailID\": \"email@example.com\"}]" > students.json
```

7. Run the binary

```
    $ ./CheckMyResult
```

## Contribute

You are always welcome to open an issue or provide a pull-request!

## License

Built under [MIT](LICENSE) license.
