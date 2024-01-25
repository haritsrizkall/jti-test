ALTER TABLE users
ADD CONSTRAINT `users_email_unique` UNIQUE (`email`);

ALTER TABLE phones
ADD CONSTRAINT `phones_number_unique` UNIQUE (`number`);
