module.exports = function(context, callback) {
    const request = context.request;
    const body = request.body;
    const payload = body.payload;
    const result = payload.result;
    const misbehaving = Object.keys(result).filter(
        (k) => {
            const o = JSON.parse(result[k]);
            const past = Date.now() - 5 * 60 * 1000
            return o.status === 'error'
                || o.battery < 5
                || o.atime < past;
        }
    );
    payload.misbehaving = misbehaving;
    payload.msg = 'Some node is not working properly: '+JSON.stringify(misbehaving);

    if (misbehaving.length === 0) {
        callback(200, null);
        return;
    }
    callback(200, payload);
};
