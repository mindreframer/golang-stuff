echo ""
echo "Unit Tests"

if [ "$1" = 'valgrind' ]
then
    VALGRIND_CMD="valgrind --leak-check=full"
    VALGRIND="$VALGRIND_CMD --log-file=/tmp/valgrind.log"
fi

# Loop over compiled tests and run them.
for test_file in tests/*_tests
do
    # Only execute if result is a file.
    if test -f $test_file
    then
        # Clear tmp directory.
        rm -rf tmp
        mkdir -p tmp

        # Log execution to file.
        if $VALGRIND ./$test_file 2>&1 > /tmp/sky-test.log
        then
            rm -f /tmp/sky-test.log

            # Check valgrind log if enabled.
            if [ -n "$VALGRIND" ]; then
                VALGRIND_ERROR_SUMMARY=`grep "ERROR SUMMARY: 0 errors from 0 contexts" /tmp/valgrind.log`
                
                if [ -z "$VALGRIND_ERROR_SUMMARY" ]; then
                    cat /tmp/valgrind.log
                    echo ""
                    echo "Run the following to reproduce:"
                    echo ""
                    echo "  $VALGRIND_CMD ./$test_file"
                    echo ""
                    exit 1
                fi
            fi
        else
            # If error occurred then print off log.
            cat /tmp/sky-test.log
            exit 1
        fi
    fi
done

echo ""
