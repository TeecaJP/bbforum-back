-- サンプルデータ挿入スクリプト

-- チーム（ボード）の作成
INSERT INTO boards (id, name, created_at) VALUES
('buffaloes-board', 'Orix Buffaloes', NOW()),
('marines-board', 'Chiba Lotte Marines', NOW()),
('hawks-board', 'Fukuoka SoftBank Hawks', NOW()),
('eagles-board', 'Tohoku Rakuten Golden Eagles', NOW()),
('lions-board', 'Saitama Seibu Lions', NOW()),
('fighters-board', 'Hokkaido Nippon-Ham Fighters', NOW()),
('giants-board', 'Yomiuri Giants', NOW()),
('tigers-board', 'Hanshin Tigers', NOW()),
('dragons-board', 'Chunichi Dragons', NOW()),
('swallows-board', 'Tokyo Yakult Swallows', NOW()),
('carp-board', 'Hiroshima Toyo Carp', NOW()),
('baystars-board', 'Yokohama DeNA BayStars', NOW());
SHOW TABLES;

-- サンプルユーザーの作成
INSERT INTO users (id, email, name, image, type, created_at, updated_at) VALUES
('user1', 'user1@example.com', 'Baseball Fan 1', 'https://example.com/user1.jpg', 'regular', NOW(), NOW()),
('user2', 'user2@example.com', 'Baseball Expert', 'https://example.com/user2.jpg', 'expert', NOW(), NOW()),
('user3', 'user3@example.com', 'Team Supporter', 'https://example.com/user3.jpg', 'regular', NOW(), NOW()),
('user4', 'user4@example.com', 'Sports Journalist', 'https://example.com/user4.jpg', 'journalist', NOW(), NOW()),
('user5', 'user5@example.com', 'Casual Viewer', 'https://example.com/user5.jpg', 'regular', NOW(), NOW());
SHOW TABLES;

-- 各ボードに投稿を追加
INSERT INTO posts (id, board_id, user_id, content, created_at) VALUES
('post1', 'buffaloes-board', 'user1', 'Excited for the Buffaloes game tonight!', NOW()),
('post2', 'marines-board', 'user2', 'Marines pitcher looked great in yesterdays game.', NOW()),
('post3', 'hawks-board', 'user3', 'Hawks are on fire this season!', NOW()),
('post4', 'eagles-board', 'user4', 'Eagles need to improve their batting average.', NOW()),
('post5', 'lions-board', 'user5', 'First time at a Lions game, the atmosphere is amazing!', NOW()),
('post6', 'fighters-board', 'user1', 'Fighters new stadium is impressive.', NOW()),
('post7', 'giants-board', 'user2', 'Giants vs Tigers rivalry never gets old.', NOW()),
('post8', 'tigers-board', 'user3', 'Tigers fans are the most passionate!', NOW()),
('post9', 'dragons-board', 'user4', 'Dragons pitching rotation is looking strong this year.', NOW()),
('post10', 'swallows-board', 'user5', 'Swallows game was rained out, disappointing.', NOW()),
('post11', 'carp-board', 'user1', 'Carps home run derby was exciting to watch.', NOW()),
('post12', 'baystars-board', 'user2', 'BayStars are underdogs but theyre my favorite team.', NOW());

-- タグの追加a
INSERT INTO tags (id, name) VALUES
('tag1', 'GameDay'),
('tag2', 'BaseballLife'),
('tag3', 'FanZone'),
('tag4', 'PitchingAce'),
('tag5', 'HomeRunHero');

-- 投稿にタグを関連付け
INSERT INTO post_tags (post_id, tag_id) VALUES
('post1', 'tag1'), ('post1', 'tag3'),
('post2', 'tag2'), ('post2', 'tag4'),
('post3', 'tag1'), ('post3', 'tag5'),
('post4', 'tag2'), ('post4', 'tag4'),
('post5', 'tag3'), ('post5', 'tag1'),
('post6', 'tag2'), ('post6', 'tag3'),
('post7', 'tag1'), ('post7', 'tag5'),
('post8', 'tag3'), ('post8', 'tag2'),
('post9', 'tag4'), ('post9', 'tag2'),
('post10', 'tag1'), ('post10', 'tag3'),
('post11', 'tag5'), ('post11', 'tag2'),
('post12', 'tag3'), ('post12', 'tag1');

-- 返信の追加
INSERT INTO posts (id, board_id, user_id, content, reply_to, created_at) VALUES
('reply1', 'buffaloes-board', 'user2', 'Me too! Who do you think will be the starting pitcher?', 'post1', NOW()),
('reply2', 'marines-board', 'user3', 'Agreed! His fastball was clocking at 98 mph!', 'post2', NOW()),
('reply3', 'hawks-board', 'user4', 'Their batting lineup is unstoppable this year.', 'post3', NOW()),
('reply4', 'eagles-board', 'user5', 'Theyve been practicing hard, Im sure theyll improve soon.', 'post4', NOW()),
('reply5', 'lions-board', 'user1', 'Wait until you see them in a playoff game!', 'post5', NOW());