insert into accounts values (null, 1, "uuid-gen-12345", "reseller", "My Name", "firstuser",md5("password"),"firstuser@email.org",now());
insert into hotspots values (null,1,"HS Test","HS for development",now());
insert into accounts values (null, 1, "uuid-gen-99999", "customer", "Customer Account", "cust1",md5("cust1"),"cust1@example.com",now());
insert into accounts_hotspots values (null, 3, 1);
