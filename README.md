# Sorting & Searching Lab

A small, dependency-free teaching app: a Go HTTP microservice generates algorithm traces, and a browser UI visualizes them one step at a time.

## Run it

Requires Go 1.22 or newer.

```sh
go run .
```

Open [http://localhost:8080](http://localhost:8080). Set `PORT` to change the listening port.

## What’s included

- Sorting: bubble, selection, insertion, merge, and quick sort.
- Searching: linear and binary search. Binary search requires ascending input.
- Sample datasets, editable comma-separated integers, playback controls, and step-by-step comparison/swap counts.
- JSON endpoints: `GET /api/algorithms`, `POST /api/sort`, and `POST /api/search`.

Example request:

```sh
curl -X POST http://localhost:8080/api/sort \
  -H 'Content-Type: application/json' \
  -d '{"algorithm":"merge","values":[8,3,5,1]}'
```

The sort and search endpoints accept between 1 and 100 integer values. The visual interface caps input at 30 values to keep traces readable.

## Test

```sh
go test ./...
```
