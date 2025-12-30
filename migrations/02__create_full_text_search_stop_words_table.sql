CREATE TABLE full_text_search_stop_words (
    id INT AUTO_INCREMENT PRIMARY KEY,
    word VARCHAR(255) NOT NULL
) ENGINE=InnoDB;

INSERT INTO full_text_search_stop_words (word) VALUES 
('з'), ('зі'), ('із'), ('без'), ('на'), ('в'), ('у'), ('для'), ('від'), ('до'), ('за'), ('по'),
('і'), ('та'), ('й'), ('ніж'), ('не'), ('тільки'), ('лише'), ('навіть'), ('ні'), ('же'), ('ж');

 