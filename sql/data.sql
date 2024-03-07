INSERT INTO users (nome, nick, email, password)
values
("user 1", "user_1", "user1@gmail.com", "$2a$10$BpvIAF83paKczAQjjwBH8Om1h/0Wrk3arODnlQTHm8lly./IqRaIa"),
("user 2", "user_2", "user2@gmail.com", "$2a$10$BpvIAF83paKczAQjjwBH8Om1h/0Wrk3arODnlQTHm8lly./IqRaIa"),
("user 3", "user_3", "user3@gmail.com", "$2a$10$BpvIAF83paKczAQjjwBH8Om1h/0Wrk3arODnlQTHm8lly./IqRaIa");

INSERT INTO followers(user_id, follower_id)
values
(1, 2)
(3, 1)
(1, 3)