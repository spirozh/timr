#!/bin/sh

cat > /tmp/kill.sh <<END
#!/bin/sh
id=\$(lsof -i :8080 | tail -1 | tr -s ' ' | cut -d ' ' -f 2 )
[[ -z "\$id" ]] || kill -9 \$id  
END
chmod +x /tmp/kill.sh

cat > /tmp/run.sh <<END
#!/bin/sh

/tmp/kill.sh
go run ./cmd/timrd
END
chmod +x /tmp/run.sh

reflex --start-service -- /tmp/run.sh

/tmp/kill.sh
rm /tmp/run.sh
rm /tmp/kill.sh
