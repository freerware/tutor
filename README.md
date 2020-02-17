# tutor
An HTTP API demonstrating example usage of various libraries created, owned,
and maintained by freerware.

## Quickstart

**(! NOTE: REQUIRES `docker` !)**

Fire it up! 🔥

```bash
make
```

Fire it up container-less! 🔥

```bash
make local
```

Tear it down! 🚧

```bash
make down 
```

Debug! 🔍🐛

```bash
make debug
```

Debug the database! 🔍🐛

```bash
make debug-db
```

Build! 🔨

```bash
make bins
```

Clean! 🧽

```bash
make clean 
```

## cURL Examples

Create a new `account`:
```bash
cd ./curl/account/ && curl -K post_account.curl http://127.0.0.1:8000/accounts/ && cd ../../
```

Upsert an`account`:
```bash
cd ./curl/account/ && curl -K put_account.curl http://127.0.0.1:8000/accounts/04b8db89-cf81-47c8-ae26-b48ae60f1e09 && cd ../../
```

Retrieve an existing `account`:
```bash
cd ./curl/account/ && curl -K get_account.curl http://127.0.0.1:8000/accounts/04b8db89-cf81-47c8-ae26-b48ae60f1e09 && cd ../../
```

Remove an existing `account`:
```bash
cd ./curl/account/ && curl -K delete_account.curl http://127.0.0.1:8000/accounts/04b8db89-cf81-47c8-ae26-b48ae60f1e09 && cd ../../
```
