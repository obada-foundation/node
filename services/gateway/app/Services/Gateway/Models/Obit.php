<?php

declare(strict_types=1);

namespace App\Services\Gateway\Models;

use Illuminate\Database\Eloquent\Model;
use OwenIt\Auditing\Contracts\Auditable;
use OwenIt\Auditing\Auditable as AuditableConcern;

class Obit extends Model implements Auditable {

    use AuditableConcern;

    const FUNCTIONAL_STATUS = 'FUNCTIONAL';

    const NON_FUNCTIONAL_STATUS = 'NON_FUNCTIONAL';

    const DISPOSED_STATUS = 'DISPOSED';

    const STOLEN_STATUS = 'STOLEN';

    const DISABLED_BY_OWNER_STATUS = 'DISABLED_BY_OWNER';

    const STATUSES = [
        self::FUNCTIONAL_STATUS,
        self::NON_FUNCTIONAL_STATUS,
        self::DISPOSED_STATUS,
        self::STOLEN_STATUS,
        self::DISABLED_BY_OWNER_STATUS
    ];

    protected $table = 'gateway_view';

    protected $casts = [
        'obit_did_versions' => 'array',
        'metadata'          => 'array',
        'doc_links'         => 'array',
        'structured_data'   => 'array',
        'modified_at'       => 'datetime'
    ];

    protected $guarded = [];

    public $timestamps = false;
}