-- Hapus kolom author_name
ALTER TABLE comments
DROP COLUMN author_name;

-- Tambahkan kolom author_id dan tetapkan foreign key
ALTER TABLE comments
ADD COLUMN author_id INT,
ADD CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES users(id);
