INSERT INTO
  users (email, password_hash, name, role)
VALUES
  (
    'admin@example.com',
    '$2y$10$TrEEPSHI1lj9RyCoH8Sv8.aIrUl0I0MABelSeDHdSIKmWAJxyM5dO',
    'Admin User',
    'admin'
  ),
  (
    'somchai@example.com',
    '$2y$10$pcAswHaW9kkayCXIIyLn6uEHpd1yBcpmaVW.avbkIjHubo/AgWQrW',
    'สมชาย ใจดี',
    'staff'
  ),
  (
    'malee@example.com',
    '$2y$10$pcAswHaW9kkayCXIIyLn6uEHpd1yBcpmaVW.avbkIjHubo/AgWQrW',
    'มาลี รักงาน',
    'staff'
  );

INSERT INTO
  equipments (code, name, category, status)
VALUES
  (
    'NB-001',
    'Dell Latitude 5440',
    'notebook',
    'available'
  ),
  (
    'NB-002',
    'MacBook Pro 14"',
    'notebook',
    'available'
  ),
  (
    'NB-003',
    'Lenovo ThinkPad T14',
    'notebook',
    'borrowed'
  ),
  (
    'PRJ-001',
    'Epson EB-982W',
    'projector',
    'available'
  ),
  ('CAM-001', 'Canon EOS R50', 'camera', 'repair'),
  (
    'MON-001',
    'Dell UltraSharp U2723QE',
    'monitor',
    'available'
  );