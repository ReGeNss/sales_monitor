CREATE TABLE full_text_search_stop_words (
    value VARCHAR(255) NOT NULL
) ENGINE=InnoDB;

INSERT INTO full_text_search_stop_words (value) VALUES 
('з'), ('зі'), ('із'), ('без'), ('на'), ('в'), ('у'), ('для'), ('від'), ('до'), ('за'), ('по'),
('і'), ('та'), ('й'), ('ніж'), ('не'), ('тільки'), ('лише'), ('навіть'), ('ні'), ('же'), ('ж'),
('смак'), ('смаком');

 