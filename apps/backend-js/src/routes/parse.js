import express from 'express';
import multer from 'multer';

const router = express.Router();
const upload = multer({ limits: { fileSize: 10 * 1024 * 1024 } });

export function createParseRouter(parserService) {
  router.post('/delegation', async (req, res) => {
    try {
      const { token } = req.body;
      if (!token) {
        return res.status(400).json({ error: 'Token is required' });
      }

      const tokenBytes = Buffer.from(token, 'base64');
      const result = await parserService.parseDelegation(tokenBytes);
      res.json(result);
    } catch (error) {
      res.status(422).json({ 
        error: 'Failed to parse delegation',
        message: error.message 
      });
    }
  });

  router.post('/chain', async (req, res) => {
    try {
      const { token } = req.body;
      if (!token) {
        return res.status(400).json({ error: 'Token is required' });
      }

      const tokenBytes = Buffer.from(token, 'base64');
      const result = await parserService.parseDelegationChain(tokenBytes);
      res.json(result);
    } catch (error) {
      res.status(422).json({ 
        error: 'Failed to parse chain',
        message: error.message 
      });
    }
  });

  router.post('/invocation', async (req, res) => {
    try {
      const { token } = req.body;
      if (!token) {
        return res.status(400).json({ error: 'Token is required' });
      }

      const tokenBytes = Buffer.from(token, 'base64');
      const result = await parserService.parseInvocation(tokenBytes);
      res.json(result);
    } catch (error) {
      res.status(422).json({ 
        error: 'Failed to parse invocation',
        message: error.message 
      });
    }
  });

  router.post('/delegation/file', upload.single('file'), async (req, res) => {
    try {
      if (!req.file) {
        return res.status(400).json({ error: 'File is required' });
      }

      const result = await parserService.parseDelegation(req.file.buffer);
      res.json(result);
    } catch (error) {
      res.status(422).json({ 
        error: 'Failed to parse file',
        message: error.message 
      });
    }
  });

  return router;
}
