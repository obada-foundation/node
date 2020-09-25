<?php

declare(strict_types=1);

namespace App\Http\Requests\Obit;

use Pearl\RequestValidate\RequestAbstract;

class CreateRequest extends RequestAbstract {

    public function authorize() {
        return true;
    }

    public function rules() {
        return [
            'obit_did'    => 'required',
            'usn'         => 'required',
            'modified_at' => 'required|date'
        ];
    }
}