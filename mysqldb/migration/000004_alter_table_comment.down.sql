-- Hapus foreign key constraint dan kolom author_id
ALTER TABLE comments
DROP FOREIGN KEY fk_author_id,
DROP COLUMN author_id;

-- Tambahkan kembali kolom author_name
ALTER TABLE comments
ADD COLUMN author_name VARCHAR(255) NOT NULL;
