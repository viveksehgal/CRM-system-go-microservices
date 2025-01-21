brew services start postgresql@14

createdb auth_service_crm

psql auth_service_crm 
CREATE USER postgres  WITH PASSWORD '12345';

