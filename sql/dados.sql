INSERT INTO usuarios( nome, nick, email, senha)
values
("user 1", "user1", "user1@gmail.com", "123456"),
("user 2", "user2", "user2@gmail.com", "12456"),
("user 3", "user3", "user3@gmail.com", "123455");

INSERT INTO seguidores(usuario_id, seguidor_id)
values
(1,2),
(2,1),
(3,1),
(2,3);