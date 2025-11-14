DELIMITER $$

CREATE PROCEDURE `insert_otp`(
    IN p_phone CHAR(11),
    IN p_expires_at TIMESTAMP
)
BEGIN
    DECLARE v_code CHAR(5);

    -- تولید کد تصادفی ۵ رقمی
    SET v_code = LPAD(FLOOR(RAND() * 100000), 5, '0');

    -- درج یا جایگزینی رکورد (بدون غیرفعال سازی بقیه)
    INSERT INTO otp_code (phone, code, expire_at)
    VALUES (p_phone, v_code, p_expires_at)
    ON DUPLICATE KEY UPDATE code = v_code, create_at = NOW(), expire_at = p_expires_at;

    -- برگرداندن کد و شماره تلفن
    SELECT v_code AS code, p_phone AS phone;
END$$

DELIMITER ;



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
      AND is_used = FALSE;

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



DELIMITER $$

CREATE PROCEDURE `insert_refresh_token`(
    IN p_user_id MEDIUMINT UNSIGNED,
    IN p_token CHAR(255),
    IN p_expires_at DATETIME
)
BEGIN
    -- توجه به اینکه جدول باید به نام refresh_tokens باشد
    INSERT INTO refresh_tokens(user_id, token, expires_at)
    VALUES (p_user_id, p_token, p_expires_at)
    ON DUPLICATE KEY UPDATE
        token = p_token,
        expires_at = p_expires_at,
        revoked = FALSE,
        created_at = NOW();
END$$

DELIMITER ;



DELIMITER $$

CREATE PROCEDURE `validate_refresh_token`(
    IN p_token CHAR(255)
)
BEGIN
    DECLARE v_user_id MEDIUMINT UNSIGNED;

    -- استفاده از NOT FOUND به جای NULL برای بررسی عدم یافتن رکورد
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_user_id = NULL;

    SELECT user_id INTO v_user_id
    FROM refresh_tokens
    WHERE token = p_token AND revoked = FALSE AND expires_at > NOW();

    IF v_user_id IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Refresh token نامعتبر یا منقضی شده';
    END IF;

    SELECT v_user_id AS user_id;
END$$

DELIMITER ;
