UPDATE students 
SET a01 = CASE 
           WHEN toan IS NOT NULL AND ly IS NOT NULL AND anh IS NOT NULL THEN toan + ly + anh
           ELSE NULL
         END;

UPDATE students 
SET a00 = CASE 
           WHEN toan IS NOT NULL AND ly IS NOT NULL AND hoa IS NOT NULL THEN toan + ly + hoa
           ELSE NULL
         END;

UPDATE students 
SET d01 = CASE 
           WHEN toan IS NOT NULL AND anh IS NOT NULL AND van IS NOT NULL THEN toan + van + anh
           ELSE NULL
         END;
