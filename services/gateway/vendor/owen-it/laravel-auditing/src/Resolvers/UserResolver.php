<?php

namespace OwenIt\Auditing\Resolvers;

use App\User;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Config;

class UserResolver implements \OwenIt\Auditing\Contracts\UserResolver
{
    /**
     * {@inheritdoc}
     */
    public static function resolve()
    {
        return new User();
    }
}
