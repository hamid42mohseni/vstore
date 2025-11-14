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
