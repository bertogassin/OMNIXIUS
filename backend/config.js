export const config = {
  port: process.env.PORT || 3000,
  jwtSecret: process.env.JWT_SECRET || 'omnixius-dev-secret-change-in-production',
  jwtExpires: process.env.JWT_EXPIRES || '7d',
  bcryptRounds: 12,
  maxLoginAttempts: 5,
  loginWindowMs: 15 * 60 * 1000, // 15 min
  uploadDir: 'uploads',
  allowedImageTypes: ['image/jpeg', 'image/png', 'image/webp'],
  maxFileSize: 5 * 1024 * 1024, // 5 MB
};
