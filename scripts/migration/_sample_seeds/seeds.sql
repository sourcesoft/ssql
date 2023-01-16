INSERT INTO public."user" (id,username,email,email_verified,active,created_at,updated_at,deleted_at) VALUES
	 ('2d9392f9-c7ab-45e6-8a9e-c883ad4460c9','test','test@test.com',true,true,1670384338,1670384338,NULL);
INSERT INTO public."user" (id,username,email,email_verified,active,created_at,updated_at,deleted_at) VALUES
	 ('7d83a8c0-084c-4954-92d7-87bf4b8256fe','test2','test2@test.com',true,true,1670384349,1670384349,NULL);
INSERT INTO public."user" (id,username,email,email_verified,active,created_at,updated_at,deleted_at) VALUES
	 ('6602363b-0357-46cd-8802-345f26a99064','test3','test3@test.com',true,true,1670384359,1670384359,NULL);
INSERT INTO public."user" (id,username,email,email_verified,active,created_at,updated_at,deleted_at) VALUES
	 ('860bc661-47ff-469d-9847-12c1c4a4630f','test4','test4@test.com',false,false,1670384444,1670384444,NULL);

INSERT INTO public."identity" (id,user_id,name,given_name, nick_name, family_name) VALUES
	 ('4370176f-4213-4ac6-8452-260a81c58f33','2d9392f9-c7ab-45e6-8a9e-c883ad4460c9','test','pouya', 'po', 'sanooei');
INSERT INTO public."identity" (id,user_id,name,given_name, nick_name) VALUES
	 ('761fa3af-1d3a-42f6-a10f-322493a57bf3','7d83a8c0-084c-4954-92d7-87bf4b8256fe','test2','john', 'joe');

INSERT INTO public."extras" (id, user_id, locale, admin_notes, bio, about) VALUES
	 ('c7d8954e-499d-4970-92c7-7b0375061e22','2d9392f9-c7ab-45e6-8a9e-c883ad4460c9','en','note test', 'bio test', 'about me');
INSERT INTO public."extras" (id, user_id, locale, admin_notes, bio, about) VALUES
	 ('34c04747-6f5c-4a18-9574-fb18a9ec45ee','7d83a8c0-084c-4954-92d7-87bf4b8256fe','fr','note test 2', 'bio test 2');
