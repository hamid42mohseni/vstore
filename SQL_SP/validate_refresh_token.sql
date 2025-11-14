DELIMITER $$

CREATE PROCEDURE `validate_refresh_token`(
    IN p_token CHAR(255)
)
BEGIN
    DECLARE v_user_id MEDIUMINT UNSIGNED;

    SELECT user_id INTO v_user_id
    FROM refresh_tokens
    WHERE token = p_token AND revoked = FALSE AND expires_at > NOW();

    IF v_user_id IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Refresh Token Expired';
    END IF;

    SELECT v_user_id AS user_id;
END$$

DELIMITER ;

