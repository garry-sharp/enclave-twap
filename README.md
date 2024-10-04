# API Tool

## Code Structure

```
├── README.md
├── app.log
├── go.mod
├── go.sum
├── main.go
└── pkg
    ├── api
    │   ├── api.go -- Library like calls to interact with the API
    │   ├── api_test.go
    │   ├── auth.go -- Handles authentication headers
    │   ├── auth_test.go
    │   ├── endpoints.go -- Actual API endpoints and wraps their responses
    │   ├── endpoints_test.go
    │   ├── helper.go -- Used for parsing dateTime strings in json
    │   ├── net.go -- Configuration and credential loading
    │   ├── net_test.go
    │   ├── rate.go -- Unused rate limiter
    │   ├── response_types.go -- JSON structs of responses
    │   └── types.go -- Other types used in this implementation
    ├── cli
    │   └── cobra.go -- Handles the initial CLI load on launch
    ├── logger
    │   └── logger.go -- Handles logging
    └── twap
        ├── helper.go -- A series of helper functions used in the TWAP implementation
        ├── helper_test.go
        └── twap.go -- The core TWAP implementation code
```

## Run

The code is not dockerised and runs on Golang `1.23.1`. There are two options for running it. Either settings values in a `.env` file or passing values as command line args, see [examples](#examples) for more info.

### Build

The code can be built with `go build -o assessment .` and then run by calling `./assessment`

### Examples

#### Example 1

```bash
echo "API_KEY=<YOUR_API_KEY>" > .env
echo "API_SECRET=<YOUR_API_SECRET>" >> .env
echo "BASE_URL=https://api-sandbox.enclave.market" >> .env
echo "TRADE_SIDE=buy" >> .env
echo "AMOUNT=100" >> .env
echo "DURATION=60s" >> .env
echo "MARKET=AVAX-USDC" >> .env
echo "INTERVAL=5s" >> .env
go run main.go
```

#### Example 2

```bash
go run main.go --side buy twap --duration "1m" --interval "5s" --amount "100" --market "AVAX-USDC"
```

### CLI

I decided to use cobra due to how well it's been tested to handle the CLI, additionally a .env file has been added to handle unit tests. The two work in sync with one another as to allow the CLI to be lightweight.

### Logger

A wrapper of the native golang logger that allows for different levels of logging as well as to multiple writers. This could be extended further to allow for pushing of data to external services like kafka (not in scope for this implementation)

### API Package

#### Net

This is used for actual API calls and returns the data as is, parsed into JSON structs See [Response Types](#response-types). For efficiency and following golang convention, the struct to be parsed to is passed in as a parameter. These functions are private

#### API

These functions are convenient wrappers around the [Net](#net) calls, they are intended to be used for the twap and potentially with 3rd party libraries

#### Auth

Used by the [Net](#net) calls to add the required authentication headers to the outgoing API requests.

#### Response Types

JSON structs that can be parsed and returned. There is a Generic type of `APIResponse` which contains any result type `T`. This makes sense as all responses from the API follow the convention of

---

| Name       | Type     | Required |
| ---------- | -------- | -------- |
| success    | `bool`   | `false`  |
| result     | `any`    | `true`   |
| error      | `string` | `true`   |
| error_code | `string` | `true`   |

---

### TWAP

There core of the TWAP code is here. The basic execution flow is thus

1. Do a quick sanity check on the parameters passed in.
2. Load in the API keys
3. Verify the user is authenticated //TODO can they actually make an order
4. Verify the market exists and get the increments
5. Reduce the quantity to the nearest increment (round down)
6. Check there is enough balance to perform the TWAP
7. Get the number of iterations in the TWAP and the quantities spread out as evenly as possible so each value differs my at _MOST_ the `increment` value
8. If no errors so far, then proceed to the actual TWAP execution.
    1. Create a ticker from the `time` package to send a signal on a channel every interval.
    2. Create a cancelable `context` in case there are errors mid flight.
    3. For each iteration (except the first). Wait for either a `ticker` signal or a `cancel` signal
        - if `cancel` then decrement the wait group and continue
        - if `ticker` then launch a goroutine to execute the order
        - if any goroutine fails 3 times consecutively, then trigger `cancel`

## Error Handling.

The format of error handling for this is to try and catch all possible errors before executing the main twap function. There are many cases where this may fail, such as `insufficient_funds` or not having the proper credentials.

There are a few cases where the TWAP may fail in its execution. Namely, if API permissions are revoked mid execution or the balance becomes insufficient to continue. In these cases there are no corrective actions to be taken and so each order will be attempted to be placed a maximum of `3` times before cancelling all other go routines.

## Other

-   Minimum interval size is 500ms
-   Maximum number of intervals is 1000
-   Intervals must divide perfectly into the duration

### Rounding Errors

In the code I use `big.Float` for accuracy, moving between it and `string` types where appropriate. However, there are minimum limits imposed by `Enclave` which cause issues. Namely, the `minimum increment`. If I try and buy 1 USD worth of AVAX 60 times in 1 minute I can get at most `0.001` AVAX back. This will result in all the orders being made being cancelled. This is a context issue and hard to protect against programmatically.

### Rate limiting

I did start to add a rate limiter but later removed it as I thought it out of scope. This would be added using the `golang.org/x/time/rate` package which implements a `token bucket` type limiter like the `Enclave` API uses. You would initialize with the `refill` and `burst` rates and call `Wait()` on it. I didn't want to overly complicate the project.

### Latency

It takes around 4 seconds to send a trade to the server with the REST API and read the result, this is why the ticker controls the start time of the trade and the trade confirms a few seconds later.

### Keys

We can check if the user is authenticated rather easily, however it is hard to check if the user has the appropriate permissions on the keys. This will only happen later once the TWAP actually starts

## Out of Scope

-   A rate limiting mechanism. If this were to be built in then you could simply create a rate limiter struct and make the net calls keep track of the number of requests made within a time period. But for this it's out of scope due to time limitations.
-   External logging to systems like `Kafka`.
-   Prompts to correct incorrectly set parameters.
-   Recovery in case of failure during execution. We just stop.
-   Databases, there's no persistent storage of trades after they've happened in a DB, which you would usually do for analytics or audit purposes.
-   Dry runs, nice to have but let's not over complicate it;
-   Websockets. Again this would make more sense for doing higher volumes over a shorter period but is out of scope for this project.
