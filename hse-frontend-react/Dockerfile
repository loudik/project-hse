# ---- Build stage ----
FROM node:20-alpine AS builder
WORKDIR /app

# Passed in from docker-compose.yml (build.args), sourced from the root
# .env - so these only need to be set in ONE place, not duplicated in
# hse-frontend-react/.env as well when building via Docker.
ARG VITE_MS_CLIENT_ID
ARG VITE_MS_TENANT_ID
ARG VITE_MS_REDIRECT_URI
ENV VITE_MS_CLIENT_ID=$VITE_MS_CLIENT_ID
ENV VITE_MS_TENANT_ID=$VITE_MS_TENANT_ID
ENV VITE_MS_REDIRECT_URI=$VITE_MS_REDIRECT_URI

COPY package*.json ./
RUN npm install --legacy-peer-deps
COPY . .
RUN npm run build

# ---- Serve stage ----
FROM nginx:1.27-alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
