# Redis Notes

## Some Data Structures That Redis Provides

- Hash
- List
- Set
- Sorted set
- Bitmap
- HyperLogLog
- Geospatial indexes
- Streams

## What These Data Structures Can Be Used For

- Realtime chat
- Auth session store
- Message buffers
- Media streaming
- Gaming leaderboards
- Realtime analytics

## Atomicity

Every operation in Redis is **atomic**.

- When a command is executing, Redis does not context switch and start executing another command.

## Data Is Stored In-Memory

Hence, the most common use of Redis is for caching.

```text
                <----> [MySQL]
User ---> [API]
                <----> [Redis]
```

Redis also provides configurable persistence:

- Periodically dumping data to disk
- Write-ahead logging of all commands. Every update command is logged in the append-only file (AOF), which can be used to reconstruct Redis.
- No persistence at all

## Other Key Features

- Transactions
- Pub/Sub
- TTL on keys, for anything temporary
- LRU eviction, a key eviction strategy. When the cache is full, Redis can evict keys and still accept writes.

## Concurrent Programming Models (Single Process)

Concurrent programming is all about doing more than one thing at the same time.

### Multi-Threading

Each incoming request over the network is accepted by the server and executed in a separate thread.

```text
R1 --> Increment k --> T1
R2 --> Increment k --> T2
```

How do we ensure data correctness?

- Make one thread wait while another is executing.
- Use mutexes and semaphores (pessimistic locking).

```text
Acquire LOCK
    k++
Release LOCK
```

### I/O Multiplexing (Apparent Concurrency)

- This is how event loops are implemented.

I/O system calls are blocking:

- `read()` blocks until the other side sends data.

Hence, we cannot just invoke `read()` on a socket unless we know that the other party will definitely send data.

Handling each connection on a separate thread is popular because while one thread is blocked, another thread executes. But we then need mutexes and semaphores to protect the critical section.

We need something that:

- does not keep us waiting for I/O
- can notify us when there is some movement

Core idea: use I/O monitoring calls to monitor sockets and call `read()` on the ones that have data.

```text
Single thread -->

------ }
------ }    Accept incoming TCP connection
------ }
------ ---> Read from one socket
------ }
------ }    Execute the incoming commands
------ ---> Read from one socket
------ }
------ }    Execute the incoming commands
------ }    Accept incoming TCP connection
------ ---> Read from one socket
------ }
------ }    Execute the incoming commands
```

**There is no separate process.**

**There is no separate thread.**

### What Redis Exploits

- Network I/O is slow -> waiting to receive commands
- In-memory operations are fast -> upon receiving commands, Redis can execute them very quickly

This is why Redis is:

1. Single-threaded -> no need for mutexes, semaphores, and waiting
2. Doing I/O multiplexing -> handling multiple TCP connections concurrently

## Redis Serialization Protocol (RESP)

RESP supports common data types like integers, strings, and arrays, along with a way to convey errors.

Redis sends commands as arrays of strings.

```text
PUT K V -> ["PUT", "K", "V"]
```

That array is then serialized using RESP.

### RESP Description

- Redis uses RESP as a request-response protocol.
- The client sends requests in RESP.
- The server responds in RESP.
- Every data type starts with a special character.
- The data ends with `\r\n` (CRLF).

### Simple Strings

- Simple strings start with `+`.
- They are followed by the string value, for example `PONG`.
- They end with `\r\n`.
- Example: `+PONG\r\n`
- Minimal overhead because they require `n + 3` bytes in the response.

### Integers

- Integers start with `:`.
- They are followed by the integer value.
- They end with `\r\n`.
- Example: `:1729\r\n`

### Bulk Strings

- Bulk strings start with `$`.
- They are followed by the number of bytes.
- They then include `\r\n`.
- After that comes the actual string.
- They end with another `\r\n`.

Example:

```text
$4\r\nPONG\r\n
```

Why do we need bulk strings if we have simple strings?

1. Simple strings are not binary safe.
2. Simple strings cannot contain `\r\n`.
3. Bulk strings can be used to store binary data, even a PNG image.

Some special strings:

- Empty string: `$0\r\n\r\n`
- Null value: `$-1\r\n`

### Arrays

- Arrays start with `*`.
- They are followed by the number of elements.
- They then include `\r\n`.
- After that come the RESP-encoded elements.

Example:

```text
['a', 200, "cat"] -> *3\r\n $1\r\na\r\n $3\r\ncat\r\n
```

- Null arrays: `*-1\r\n`
- Empty arrays: `*0\r\n`
- Arrays can also be nested.

### Errors

- Error messages start with `-`.
- They are followed by the message.
- They end with `\r\n`.
- Example: `-Key not found\r\n`

### Key Highlights

1. RESP is human-readable.
2. Redis is simple, which means fewer bugs.
3. RESP is performant.
4. RESP uses prefixed lengths, so we know exactly how many bytes to read and process.

## Event Loops
