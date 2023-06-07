# Distributed commit log in go
Based on Distributed Services with Go book by Travis Jeffery

Write-ahead logs (WAL), transaction logs, or commit logs - same thing. Used by storage engines, file systems, VCS (git), message queues, consensus algorithms. Log here is just a change - an object that represents what changed recently. Full log is append only sequence of changes


# Definitions
* store - file that contains the data
* index - file that contains record offset and position in the file. Can be memory mapped file to increase speed
* segment - abstraction that ties store and index
* log - all segments


# todo:
* distribution
* service discovery
* raft