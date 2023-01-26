namespace go thrift

struct GetRequest {
    1: string group
    2: string key
}

struct GetResponse {
    1: binary value
}

service Memo {
    GetResponse Get(1: GetRequest req)
}