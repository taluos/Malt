# Picker

1 What is it?
2 How to use it?
3 How can I modify it?

## 1: What is it?

Picker mainly implements the load balancing algorithms such as round robin, least connection and least response time.

We plan to implement the following load balancing algorithm in this project:

- [x] Round Robin
- [x] Weighted Round Robin
- [x] Random
- [x] Weighted Random
- [x] Hash
- [x] Consistent Hash
- [x] Power of Two Choices ( P2C )

// Mabey more?

- [ ] Least Connection
- [ ] Least Active
- [ ] Least Response Time

## 2: How dose it work?

In this part, the specific load balancing algorithm is mainly implemented through the Pick() function.
It mainly takes in a slice of []nodes and returns a selected node, selectedNode, according to its rules. 
Its core pseudo-code is as follows:

```go
for each node in nodes {
    var selectedNode node
    // Sepecific load balancing algorithm
    // ...

    if selectNode is Nil || ... {
        selectedNode = node
    }
    d := selectedNode.Pick()
    return ...
}
```

## 3: How to use it?

```go

```

## 3: How can I modify it?
