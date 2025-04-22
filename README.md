# ALERTWest Interview Question 1

## Instructions

1. Hard fork this repository (git clone + new remote) to create your own private copy.
2. Add the following users as collaborators to your private repository: @ehiggins0, @jlang-dp, and @chrisschreiber.
3. Work through the tasks outlined in the Overview section.
4. Make frequent, small PRs as you complete portions of the work. Request review for each PR and we will provide feedback on your approach.
5. Document your approach, any assumptions made, and your reasoning for design decisions.

Feel free to reach out if you have any questions or need clarification on any aspect of the problem. We're here to help ensure you understand the requirements fully.

## Overview

You are working on a database that is having performance issues with an unknown cause. Provided in this repo is a server and client. The server is the executor of queries, and the client is set up to monitor the database's utilization (cpu, io, and memory). An attempt has been made to also monitor which queries are being executed, but the method is not currently effective. Your task is to improve the system by addressing these issues, first by improving the query monitoring, then by implementing query scheduling to ensure consistent performance.

## General Requirements

- The system must recover from network or client failures without losing data or significant progress
- Provide a method to run the services locally
- Provide frequent PRs as you complete work

## Use of AI Tools

While AI tools are permitted, we are looking for a good understanding of the implemented solution. You should expect follow-up questions regarding design and code decisions made throughout the problem, and inability to explain why you made certain decisions will be taken into account.

## Current Implementation

### Server (Backend Service)

The server provides the following APIs:

- `GET /queued` returns the currently queued queries, with a response like this:

```json
{
  "query": {
    "id": "550e8400-e29b-41d4-a716-446655440000" // Unique query ID
  },
  "execution": {
    "id": "a7e3f4c2-9b8d-5e6f-7c0a-1d2b3c4d5e6f", // Unique execution ID
    "timestamp": 1740000000000 // Scheduled execution timestamp (Unix time ms)
  }
}
```

- `GET /resources` returns the resource utilization of the database, updated every second:

```json
{
  "cpu": {
    "average": 50, // Percentage over the last second
    "min": 30,
    "max": 70
  },
  "io": {
    "average": 50,
    "min": 30,
    "max": 70
  },
  "memory": {
    "average": 50,
    "min": 30,
    "max": 70
  },
  "timestamp": 1740000000000 // Unix timestamp of the update
}
```

- `POST /delay` allows you to delay a specific query execution and expects a body like this:

```json
{
  "id": "a7e3f4c2-9b8d-5e6f-7c0a-1d2b3c4d5e6f", // Execution ID
  "delay": 10 // Delay in ticks, where a tick is 100ms
}
```

> [!NOTE]
> You should only need to modify the [server](server/), but feel free to take a look around lib as well.

### Client (Monitoring Service)

The client currently polls `/queued` and `/resources` every second, but does not process the response.

## Part 1: Identify Queries being Executed

### Problem

We are currently polling the `/queued` endpoint every second to get a picture of the query patterns on the server, but we're missing data - some queries are rare (or too fast) and we need to make sure that we capture everything.

### Objective

Develop a reliable mechanism to record all queries executed by the backend service.

### Requirements

- Capture every executed query, regardless of execution duration or frequency
- Ensure the system can withstand network or client failures without losing data

### Deliverables

- Modify the server and client to support your selected architecture
- Include documentation regarding architectural decisions made throughout the process

## Part 2: Optimize Query Schedule

### Problem

Now that we have a complete query execution record and associated usage metrics, we decide that we want to mitigate the resource spikes in the database. Queries are scheduled with a default 100ms delay after being queued, and the `POST /delay` endpoint allows further postponement.

### Objective

Optimize query execution scheduling to maintain a minimum and maximum CPU, IO and memory utilization within 10% of the averages, while also taking query latency into account.

### Requirements

- Develop a client-side algorithm to estimate each query's resource utilization (CPU, IO, and memory) using historical execution data and the 1-second resource metrics
- Create a client-side algorithm to determine when to execute queries, adjusting delays to achieve the +/- 10% spread from the average
- Ensure the client can be restarted and resume scheduling based on the system state.

### Deliverables

- Update the client code to schedule queries using the appropriate algorithms
- Document the system you've designed and how your approach achieves the targetted utilization

## Frequently Asked Questions

### What is lib, and what can be used / modified in it?

Lib is the implementation of our "database". It randomly selects "queries", then executes them depending on the parameters defined in [db](./lib/db.go).

All exported functions can be used in your solution. While you should not need to modify any part of lib to solve this problem, you can change the parameters in `db` to make the problem easier / harder. Your final solution should support basically any set of parameters.

### Can additional services be added?

Yes, you may add services on an as needed basis, as long as the services are open source and the entire application can be ran locally and started with a single command.

### Where can additional questions be asked?

We encourage you to frequently ask questions while working on this question. We are more than willing to clarify anything you're uncertain about, validate your solution before implementation, or just give feedback on ideas. Feel free to include questions in PRs or reach out to us via email, Slack, or Discord, and we will try to get back to you same-day.
