import express from 'express';
import multer from 'multer';

const router = express.Router();
const upload = multer({ limits: { fileSize: 10 * 1024 * 1024 } });

export function createValidateRouter(validatorService) {
  router.post('/chain', async (req, res) => {
    try {
      const { token } = req.body;
      if (!token) {
        return res.status(400).json({ error: 'Token is required' });
      }

      const tokenBytes = Buffer.from(token, 'base64');
      const result = await validatorService.validateChain(tokenBytes);
      res.json(result);
    } catch (error) {
      res.status(500).json({ 
        error: 'Validation failed',
        message: error.message 
      });
    }
  });

  router.post('/chain/file', upload.single('file'), async (req, res) => {
    try {
      if (!req.file) {
        return res.status(400).json({ error: 'File is required' });
      }

      const result = await validatorService.validateChain(req.file.buffer);
      res.json(result);
    } catch (error) {
      res.status(500).json({ 
        error: 'Validation failed',
        message: error.message 
      });
    }
  });

  return router;
}
