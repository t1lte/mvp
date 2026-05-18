FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install


COPY . .

EXPOSE 9013

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "9013"]
