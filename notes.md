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
