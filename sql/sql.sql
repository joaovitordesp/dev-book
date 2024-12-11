CREATE DATABASE IF NOT EXISTS devbook;

USE devbook;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL UNIQUE,
    nick VARCHAR(100) UNIQUE NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    criadoEm timestamp default current_timestamp()
) ENGINE=InnoDB;


DROP TABLE IF EXISTS seguidores;

CREATE TABLE seguidores (
    usuario_id int not NULL,
    Foreign KEY(usuario_id)
    References usuarios(id)
    ON DELETE CASCADE,  

    seguidor_id int not null,
    FOREIGN KEY(seguidor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    PRIMARY KEY(usuario_id, seguidor_id)  
) ENGINE=InnoDB;


CREATE TABLE posts(
    id INT AUTO_INCREMENT PRIMARY KEY,
    titulo VARCHAR(50) NOT NULL,
    conteudo VARCHAR(255) NOT NULL,
    likes int default 0
    criadoEm TIMESTAMP DEFAULT CURRENT_TIMESTAMP()

    autor_id INT NOT NULL,
    FOREIGN KEY(autor_id)
    REFERENCES usuarios(id)
    ON DELETE CASCADE,
)ENGINE=INNODB;