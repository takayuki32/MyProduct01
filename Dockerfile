FROM ruby:3.1
RUN apt-get update && apt-get install - y \
    build-essential \
    libpq-dev \
    nodejs \
    postgresql-client \
    yarn \
WORKDIR /MyProduct01
COPY Gemfile Gemfile.lock /MyProduct01/
RUN bundle install
