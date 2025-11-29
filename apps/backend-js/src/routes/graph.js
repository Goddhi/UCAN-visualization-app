import express from 'express';
import multer from 'multer';

const router = express.Router();
const upload = multer({ limits: { fileSize: 10 * 1024 * 1024 } });

export function createGraphRouter(graphService) {
  router.post('/delegation', async (req, res) => {
    try {
      const { token } = req.body;
      if (!token) {
        return res.status(400).json({ error: 'Token is required' });
      }

      const tokenBytes = Buffer.from(token, 'base64');
      const result = await graphService.generateDelegationGraph(tokenBytes);
      res.json(result);
    } catch (error) {
      res.status(422).json({ 
        error: 'Failed to generate graph',
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
      const result = await graphService.generateInvocationGraph(tokenBytes);
      res.json(result);
    } catch (error) {
      res.status(422).json({ 
        error: 'Failed to generate invocation graph',
        message: error.message 
      });
    }
  });

  return router;
}
