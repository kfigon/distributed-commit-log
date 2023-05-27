# Distributed commit log in go
Based on Distributed Services with Go book by Travis Jeffery

Write-ahead logs (WAL), transaction logs, or commit logs - same thing. Used by storage engines, file systems, VCS (git), message queues, consensus algorithms. Log here is just a change - an object that represents what changed recently. Full log is append only sequence of changes