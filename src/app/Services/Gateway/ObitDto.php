<?php

declare(strict_types=1);

namespace App\Services\Gateway;

use Spatie\DataTransferObject\DataTransferObject;
use Laravel\Lumen\Http\Request;

class ObitDto extends DataTransferObject {

    public string $obitDID;

    public string $usn;

    public $modifiedAt;

    public static function fromRequest(Request $request): self {
        return new self([
            'obitDID'    => $request->json('obit_did'),
            'usn'        => $request->json('usn'),
            'modifiedAt' => $request->json('modified_at'),
        ]);
    }
}
