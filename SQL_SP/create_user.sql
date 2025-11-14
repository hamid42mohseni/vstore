DELIMITER $$

CREATE PROCEDURE `create_user`(
    IN p_phone CHAR(11),
    IN p_code CHAR(5)
)
BEGIN
    DECLARE v_user_id MEDIUMINT UNSIGNED;
    DECLARE v_phone CHAR(11);
    DECLARE v_exists INT DEFAULT 0;

    -- بررسی صحت کد OTP
    SELECT COUNT(*) INTO v_exists
    FROM otp_code
    WHERE phone = p_phone
      AND code = p_code

    IF v_exists = 0 THEN
        SIGNAL SQLSTATE '45001' SET MESSAGE_TEXT = 'Invalid or used OTP code';
    END IF;

    -- اگر کاربر وجود نداشت ایجاد می‌شود
    IF NOT EXISTS (SELECT 1 FROM user WHERE phone = p_phone) THEN
        INSERT INTO user (phone) VALUES (p_phone);
    END IF;

    -- دریافت user_id و phone
    SELECT user_id, phone INTO v_user_id, v_phone FROM user WHERE phone = p_phone;

    -- برگرداندن user_id و phone
    SELECT v_user_id AS user_id, v_phone AS phone;

END$$

DELIMITER ;

