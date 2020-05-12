#!/bin/sh
set -e
# Make the Coverage File
echo "mode: atomic" > coverage.txt

# Make Necessary directories needed by Test (Ideally it should get created automatically but Travis is not allowing to create it using os.MkdriAll)
# I know this is insane but nothing can be done

#Start the Test
for d in $(go list ./...); do
    echo "d="$d
    if [ $(ls | grep _test.go | wc -l) -gt 0 ]; then
        go test -v -cover -covermode atomic -coverprofile coverage.out
        if [ -f coverage.out ]; then
            sed '1d;$d' coverage.out >> coverage.txt
        fi
    fi
done


