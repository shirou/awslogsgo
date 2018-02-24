# awslogsgo

A go port of [awslogs](https://github.com/jorgebastida/awslogs).

This `aslogsgo` has almost same options to awslogs, but this is very fast!

# usage

## List log groups

```
NAME:
   awslogsgo groups - list log groups

USAGE:
   awslogsgo groups [command options] [arguments...]

OPTIONS:
   --aws-access-key-id value             aws access key id [$AWS_ACCESS_KEY_ID]
   --aws-secret-access-key value         aws secret access key [$AWS_SECRET_ACCESS_KEY]
   --aws-session-token value             aws session token [$AWS_SESSION_TOKEN]
   --profile value                       aws profile [$AWS_PROFILE]
   --aws-region value                    aws region [$AWS_REGION]
   -p PREFIX, --log-group-prefix PREFIX  List only groups matching the PREFIX (default: "/")
 ```

## List log streams

```
NAME:
   awslogsgo streams - list log stream

USAGE:
   awslogsgo streams [command options] log_group_name

OPTIONS:
   --aws-access-key-id value            aws access key id [$AWS_ACCESS_KEY_ID]
   --aws-secret-access-key value        aws secret access key [$AWS_SECRET_ACCESS_KEY]
   --aws-session-token value            aws session token [$AWS_SESSION_TOKEN]
   --profile value                      aws profile [$AWS_PROFILE]
   --aws-region value                   aws region [$AWS_REGION]
   -s START, --start START              START time (default: "1h")
   -e END, --end END                    END time
   -p value, --log-stream-prefix value  List only stream matching the prefix
```

## Get logs
```
NAME:
   awslogsgo get - get log stream

USAGE:
   awslogsgo get [command options] log_group_name log_stream_name

OPTIONS:
   --aws-access-key-id value             aws access key id [$AWS_ACCESS_KEY_ID]
   --aws-secret-access-key value         aws secret access key [$AWS_SECRET_ACCESS_KEY]
   --aws-session-token value             aws session token [$AWS_SESSION_TOKEN]
   --profile value                       aws profile [$AWS_PROFILE]
   --aws-region value                    aws region [$AWS_REGION]
   -f PATTERN, --filter-pattern PATTERN  A valid CloudWatch Logs filter PATTERN to use for filtering the response. If not provided, all the events are matched.
   -w, --watch                           Query for new log lines constantly
   -G, --no-group                        Do not display group name
   -S, --no-stream                       Do not display stream name
   --timestamp                           Add creation timestamp to the output
   --ingestion-time                      Add ingestion time to the output
   -s START, --start START               START time (default: "5m")
   -e END, --end END                     END time
   --no-color                            Do not color output
```

### date

You can also specify `1h` or `3d`. If specify this, it means relative to current time.

- m, min, mins, minute, minutes
- h, hour, hours
- d, day, days
- w, week, weeks



# Bugs

- [ ] `-w` is not work

# License

BSD License (same as awslogs)
