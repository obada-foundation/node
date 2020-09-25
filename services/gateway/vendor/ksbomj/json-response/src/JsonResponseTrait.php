<?php

declare(strict_types=1);

namespace SecurityRobot;

use Illuminate\Http\{Response, JsonResponse};

trait JsonResponseTrait {
    /**
     * @var int HTTP status code
     *
     * @see http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
     */
    protected $statusCode = 200;

    /**
     * @return int
     */
    public function getStatusCode(): int {
        return $this->statusCode;
    }

    /**
     * @param int $statusCode
     *
     * @return $this
     */
    public function setStatusCode(int $statusCode) {
        $this->statusCode = $statusCode;

        return $this;
    }

    /**
     * Generates a Response with a 404 HTTP header and a given message.
     *
     * @param string $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function responseNotFound(string $message = 'Not Found!') {
        return $this->setStatusCode(404)->respondWithError($message);
    }

    /**
     * Generates a Response with a 400 HTTP header and a given message.
     *
     * @param string $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function errorWrongArgs(string $message = 'Wrong Arguments!') {
        return $this->setStatusCode(400)->respondWithError($message);
    }

    /**
     * Generates a Response with a 401 HTTP header and a given message.
     *
     * @param string $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function errorNotAuthorized(string $message = 'Not authorized!') {
        return $this->setStatusCode(401)->respondWithError($message);
    }

    /**
     * Generates a Response with a 403 HTTP header and a given message.
     *
     * @param string $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function errorForbidden(string $message = 'Forbidden!') {
        return $this->setStatusCode(403)->respondWithError($message);
    }

    /**
     * Generates a Response with a 405 HTTP header and a given message.
     *
     * @param string $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function errorMethodNotAllowed(string $message = 'HTTP Method Not Allowed!') {
        return $this->setStatusCode(405)->respondWithError($message);
    }

    /**
     * Generates a Response with a 429 HTTP header and a given message.
     *
     * @param string $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function errorTooManyRequests($message = 'Too many requests') {
        return $this->setStatusCode(429)->respondWithError($message);
    }

    /**
     * Generates a Response with a 500 HTTP header and a given message.
     *
     * @param string $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function errorInternalError(string $message = 'Internal Error!') {
        return $this->setStatusCode(500)->respondWithError($message);
    }

    /**
     * @param array $data
     * @param array $headers
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function respond($data, $headers = []) {
        return new JsonResponse($data, $this->getStatusCode(), $headers);
    }

    /**
     * Method creates the error responds structure
     *
     * @param string $message The error message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function respondWithError(string $message) {
        return $this->respond([
            'code'    => $this->getStatusCode(),
            'message' => $message,
        ]);
    }

    /**
     * @param array       $validation_errors
     * @param string|null $message
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function respondValidationErrors(array $validation_errors, string $message = null) {
        $this->setStatusCode(422);

        return $this->respond([
            'code'    => $this->getStatusCode(),
            'message' => $message,
            'errors'  => $validation_errors,
        ]);
    }

    /**
     * Returns the empty response with 204 status code
     *
     * @see https://httpstatuses.com/204
     *
     * @return \Illuminate\Http\Response
     */
    public function successRequestWithNoData() {
        return new Response(null, 204);
    }
}
