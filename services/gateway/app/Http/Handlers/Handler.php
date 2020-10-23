<?php

declare(strict_types=1);

namespace App\Http\Handlers;

use Laravel\Lumen\Routing\Controller as BaseController;
use SecurityRobot\JsonResponseTrait;

abstract class Handler extends BaseController {

    use JsonResponseTrait;

    /**
     * Execute an action on the handler.
     *
     * @param string $method
     * @param array $parameters
     * @return \Symfony\Component\HttpFoundation\Response
     */
    public function callAction($method, $parameters)
    {
        if ($method === '__invoke') {
            return call_user_func_array([$this, $method], $parameters);
        }

        throw new BadMethodCallException('Only __invoke method can be called on handler.');
    }
}
