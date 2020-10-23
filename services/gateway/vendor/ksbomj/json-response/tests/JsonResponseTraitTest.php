<?php
namespace Tests;

use Illuminate\Http\Response;
use PHPUnit\Framework\TestCase;
use \Illuminate\Http\JsonResponse;
use SecurityRobot\JsonResponseTrait;

class JsonResponseTraitTest extends TestCase {

    /**
     * @var JsonResponseTrait
     */
    protected $trait;

    public function setUp(): void
    {
        parent::setUp();

        $this->trait = $this->getMockForTrait(JsonResponseTrait::class);
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::getStatusCode
     */
    public function it_returns_200_status_code_by_default()
    {
        $this->assertEquals(200, $this->trait->getStatusCode());
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::setStatusCode
     */
    public function it_sets_status_code_correctly()
    {
        $this->trait->setStatusCode(400);

        $this->assertEquals(400, $this->trait->getStatusCode());

        $this->trait->setStatusCode(200);

        $this->assertEquals(200, $this->trait->getStatusCode());
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::respond
     */
    public function it_responds_correctly()
    {
        $data = ['data' => ['foo', 'bar']];

        $response = $this->trait->respond($data);

        $this->assertInstanceOf(JsonResponse::class, $response);
        $this->assertEquals(200, $response->getStatusCode());
        $this->assertEquals(json_encode($data), $response->content());
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::respond
     */
    public function it_responds_correctly_with_headers()
    {
        $data = ['data' => ['foo', 'bar'], ['authorization']];

        $response = $this->trait->respond($data);

        $this->assertInstanceOf(JsonResponse::class, $response);
        $this->assertEquals(200, $response->getStatusCode());
        $this->assertEquals(json_encode($data), $response->content());
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::respondWithError
     */
    public function it_verifies_that_method_respondWithError_works()
    {
        $message = 'An error';

        $data = [
            'code'    => 200,
            'message' => $message
        ];

        $response = $this->trait->respondWithError($message);

        $this->assertInstanceOf(JsonResponse::class, $response);
        $this->assertEquals(200, $response->getStatusCode());
        $this->assertEquals(json_encode($data), $response->content());
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::errorInternalError
     */
    public function it_verifies_that_method_errorInternalError_works()
    {
        $this->assertErrors('errorInternalError', 'Internal Error!', 500);
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::responseNotFound
     */
    public function it_verifies_that_method_responseNotFound_works()
    {
        $this->assertErrors('responseNotFound', 'Not Found!', 404);
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::errorWrongArgs
     */
    public function it_verifies_that_method_errorWrongArgs_works()
    {
        $this->assertErrors('errorWrongArgs', 'Wrong Arguments!', 400);
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::errorWrongArgs
     */
    public function it_verifies_that_method_errorNotAuthorized_works()
    {
        $this->assertErrors('errorNotAuthorized', 'Not authorized!', 401);
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::errorTooManyRequests
     */
    public function it_verifies_that_method_errorTooManyRequests_works()
    {
        $this->assertErrors('errorTooManyRequests', 'Too many requests', 429);
    }

    /**
     * @test
     * @covers SecurityRobot\JsonResponseTrait::errorForbidden
     */
    public function it_verifies_that_method_errorForbidden_works()
    {
        $this->assertErrors('errorForbidden', 'Forbidden!', 403);
    }

    /**
     * @param $call_method
     * @param $default_message
     * @param $http_code
     */
    public function assertErrors($call_method, $default_message, $http_code)
    {
        $data = [
            'code'    => $http_code,
            'message' => $default_message
        ];

        $response = $this->trait->{$call_method}();

        $this->assertInstanceOf(JsonResponse::class, $response);
        $this->assertEquals($http_code, $response->getStatusCode());
        $this->assertEquals(json_encode($data), $response->content());

        $message = 'Custom message';

        $data = [
            'code' => $http_code,
            'message'   => $message
        ];

        $response = $this->trait->{$call_method}($message);

        $this->assertInstanceOf(JsonResponse::class, $response);
        $this->assertEquals($http_code, $response->getStatusCode());
        $this->assertEquals(json_encode($data), $response->content());
    }

	/**
	 * @test
	 */
    public function it_tests_that_successRequestWithNoData_works() {
	    $response = $this->trait->successRequestWithNoData();

	    $this->assertInstanceOf(Response::class, $response);
	    $this->assertEquals(204, $response->getStatusCode());
	    $this->assertEquals('', $response->content());
    }
}
