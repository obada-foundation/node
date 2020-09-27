<?php

declare(strict_types=1);

namespace App\Http\Requests\Obit;

use Pearl\RequestValidate\RequestAbstract;

class UpdateRequest extends RequestAbstract {

    public function authorize() {
        return true;
    }

    public function rules() {
        return [

        ];
    }
}
