select *
from task;

SET GLOBAL general_log = 'ON';
SET GLOBAL log_output = 'TABLE';

SELECT *
FROM mysql.general_log
WHERE command_type = 'Query'
order by event_time desc;

echo 'export APP_DATABASE_USERNAME="root"' >> ~/.bashrc
echo 'export APP_DATABASE_PASSWORD="root"' >> ~/.bashrc