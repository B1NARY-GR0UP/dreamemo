namespace go thrift

struct GetRequest {
    1: string group
    2: string key
}

struct GetResponse {
    1: optional binary value
    2: optional double qps
}

service Memo {
    GetResponse Get(1: GetRequest req)
}