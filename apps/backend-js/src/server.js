import express from 'express';
import cors from 'cors';
import swaggerUi from 'swagger-ui-express';
import { ParserService } from './services/parser.js';
import { ValidatorService } from './services/validator.js';
import { GraphService } from './services/graph.js';
import { createParseRouter } from './routes/parse.js';
import { createValidateRouter } from './routes/validate.js';
import { createGraphRouter } from './routes/graph.js';
import { createUcanValidator } from './middleware/ucanvalidator.js';
import { swaggerSpec } from './swagger/config.js';
import { generateExampleTokens } from './swagger/examples.js';

const app = express();
const PORT = process.env.PORT || 8081;

// Middleware
app.use(cors());
app.use(express.json());

// Generate example tokens on startup
let exampleTokens = {};
generateExampleTokens()
  .then(tokens => {
    exampleTokens = tokens;
    console.log('✅ Generated example UCAN tokens');
    console.log('   Parse token:', tokens.parse.substring(0, 50) + '...');
    console.log('   Validate token:', tokens.validate.substring(0, 50) + '...');
    console.log('   Graph token:', tokens.graph.substring(0, 50) + '...');
    console.log('   Admin token:', tokens.admin.substring(0, 50) + '...');
  })
  .catch(err => console.warn('⚠️  Failed to generate example tokens:', err.message));

// Initialize services
const parserService = new ParserService();
const validatorService = new ValidatorService(parserService);
const graphService = new GraphService(parserService);

// UCAN validation middleware (DISABLED by default for easy testing)
app.use(createUcanValidator());

// Swagger UI with UCAN security
app.use('/swagger', swaggerUi.serve, swaggerUi.setup(swaggerSpec, {
  customCss: '.swagger-ui .topbar { display: none }',
  customSiteTitle: 'UCAN Visualization API',
  swaggerOptions: {
    persistAuthorization: true
  }
}));

// Endpoint to get example tokens
app.get('/api/examples/tokens', (req, res) => {
  res.json({
    message: 'Example UCAN tokens for testing',
    tokens: exampleTokens,
    usage: {
      parse: `Authorization: Bearer ${exampleTokens.parse}`,
      validate: `Authorization: Bearer ${exampleTokens.validate}`,
      graph: `Authorization: Bearer ${exampleTokens.graph}`,
      admin: `Authorization: Bearer ${exampleTokens.admin}`
    },
    howToUse: [
      '1. Copy one of the tokens above',
      '2. Go to /swagger',
      '3. Click "Authorize" button',
      '4. Paste: Bearer <token>',
      '5. Try the endpoints'
    ]
  });
});

// Swagger JSON
app.get('/swagger.json', (req, res) => {
  res.json(swaggerSpec);
});

// Routes
app.use('/api/parse', createParseRouter(parserService));
app.use('/api/validate', createValidateRouter(validatorService));
app.use('/api/graph', createGraphRouter(graphService));

// Health check
app.get('/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: 'ucan-backend-js',
    timestamp: new Date().toISOString(),
    version: '1.0.0'
  });
});

// Root redirect
app.get('/', (req, res) => {
  res.redirect('/swagger');
});

app.listen(PORT, () => {
  console.log(`UCAN Backend JS running on http://localhost:${PORT}`);
  console.log(`Health check: http://localhost:${PORT}/health`);
  console.log(`API Docs: http://localhost:${PORT}/swagger`);
  console.log(`Example Tokens: http://localhost:${PORT}/api/examples/tokens`);
  console.log(``);
});