DELIMITER $$

CREATE PROCEDURE `insert_refresh_token`(
    IN p_user_id MEDIUMINT UNSIGNED,
    IN p_token CHAR(255),
    IN p_expires_at DATETIME
)
BEGIN
    INSERT INTO refresh_token(user_id, token, expires_at)
    VALUES (p_user_id, p_token, p_expires_at)
    ON DUPLICATE KEY UPDATE
        token = p_token,
        expires_at = p_expires_at,
        revoked = FALSE,
        created_at = NOW();
END$$

DELIMITER ;
